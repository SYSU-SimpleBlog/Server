tensorrt fp32 fp16 tutorial with caffe pytorch minist model

# Series

- [Part 1: install and configure tensorrt 4 on ubuntu 16.04](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me%2Fpost%2Fdacc4196%2F)
- **Part 2: tensorrt fp32 fp16 tutorial**
- [Part 3: tensorrt int8 tutorial](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me%2Fpost%2F30e0cb19%2F)

# Code Example

## include headers



```cpp
#include <assert.h>
#include <sys/stat.h>
#include <time.h>

#include <iostream>
#include <fstream>
#include <sstream>
#include <iomanip>
#include <cmath>
#include <algorithm>

#include <cuda_runtime_api.h>

#include "NvCaffeParser.h"
#include "NvOnnxConfig.h"
#include "NvOnnxParser.h"
#include "NvInfer.h"
#include "common.h"

using namespace nvinfer1;
using namespace nvcaffeparser1;

static Logger gLogger;

// Attributes of MNIST Caffe model
static const int INPUT_H = 28;
static const int INPUT_W = 28;
static const int OUTPUT_SIZE = 10;
//const char* INPUT_BLOB_NAME = "data";
const char* OUTPUT_BLOB_NAME = "prob";
const std::string mnist_data_dir = "data/mnist/";


// Simple PGM (portable greyscale map) reader
void readPGMFile(const std::string& fileName, uint8_t buffer[INPUT_H * INPUT_W])
{
    readPGMFile(fileName, buffer, INPUT_H, INPUT_W);
}
```

## caffe model to tensorrt



```cpp
void caffeToTRTModel(const std::string& deployFilepath,       // Path of Caffe prototxt file
                     const std::string& modelFilepath,        // Path of Caffe model file
                     const std::vector<std::string>& outputs, // Names of network outputs
                     unsigned int maxBatchSize,               // Note: Must be at least as large as the batch we want to run with
                     IHostMemory*& trtModelStream)            // Output buffer for the TRT model
{
    // Create builder
    IBuilder* builder = createInferBuilder(gLogger);

    // Parse caffe model to populate network, then set the outputs
    std::cout << "Reading Caffe prototxt: " << deployFilepath << "\n";
    std::cout << "Reading Caffe model: " << modelFilepath << "\n";
    INetworkDefinition* network = builder->createNetwork();
    ICaffeParser* parser = createCaffeParser();

    bool useFp16 = builder->platformHasFastFp16();
    std::cout << "platformHasFastFp16: " << useFp16 << "\n";

    bool useInt8 = builder->platformHasFastInt8();
    std::cout << "platformHasFastInt8: " << useInt8 << "\n";

    // create a 16-bit model if it's natively supported
    DataType modelDataType = useFp16 ? DataType::kHALF : DataType::kFLOAT; 
    
    const IBlobNameToTensor* blobNameToTensor = parser->parse(deployFilepath.c_str(),
                                                              modelFilepath.c_str(),
                                                              *network,
                                                              modelDataType);
    // Specify output tensors of network
    // ERROR: Network must have at least one output
    for (auto& s : outputs){
        std::cout<<"output = "<< s.c_str() << std::endl;
        network->markOutput(*blobNameToTensor->find(s.c_str())); // prob
    } 

    builder->setMaxBatchSize(maxBatchSize);
    builder->setMaxWorkspaceSize(1 << 20);

    // set up the network for paired-fp16 format if available
    if(useFp16)
        builder->setFp16Mode(true);

    // Build engine
    ICudaEngine* engine = builder->buildCudaEngine(*network);
    assert(engine);

    // Destroy parser and network
    network->destroy();
    parser->destroy();

    // Serialize engine and destroy it
    trtModelStream = engine->serialize();
    engine->destroy();
    builder->destroy();

    //shutdownProtobufLibrary();
}
```

## pytorch onnx to tensorrt



```cpp
void onnxToTRTModel( const std::string& modelFilepath,        // name of the onnx model 
                     unsigned int maxBatchSize,            // batch size - NB must be at least as large as the batch we want to run with
                     IHostMemory *&trtModelStream)      // output buffer for the TensorRT model
{
    // create the builder
    IBuilder* builder = createInferBuilder(gLogger);

    nvonnxparser::IOnnxConfig* config = nvonnxparser::createONNXConfig();
    config->setModelFileName(modelFilepath.c_str());
    
    nvonnxparser::IONNXParser* parser = nvonnxparser::createONNXParser(*config);
    
    //Optional - uncomment below lines to view network layer information
    //config->setPrintLayerInfo(true);
    //parser->reportParsingInfo();
    
    if (!parser->parse(modelFilepath.c_str(), DataType::kFLOAT))
    {
        string msg("failed to parse onnx file");
        gLogger.log(nvinfer1::ILogger::Severity::kERROR, msg.c_str());
        exit(EXIT_FAILURE);
    }
    
    if (!parser->convertToTRTNetwork()) {
        string msg("ERROR, failed to convert onnx network into TRT network");
        gLogger.log(nvinfer1::ILogger::Severity::kERROR, msg.c_str());
        exit(EXIT_FAILURE);
    }
    nvinfer1::INetworkDefinition* network = parser->getTRTNetwork();
    
    // Build the engine
    builder->setMaxBatchSize(maxBatchSize);
    builder->setMaxWorkspaceSize(1 << 20);

    ICudaEngine* engine = builder->buildCudaEngine(*network);
    assert(engine);

    // we don't need the network any more, and we can destroy the parser
    network->destroy();
    parser->destroy();

    // serialize the engine, then close everything down
    trtModelStream = engine->serialize();
    engine->destroy();
    builder->destroy();

    //shutdownProtobufLibrary();
}
```

## do inference



```cpp
void doInference(IExecutionContext& context, float* input, float* output, int batchSize)
{
    const ICudaEngine& engine = context.getEngine();
    // Pointers to input and output device buffers to pass to engine.
    // Engine requires exactly IEngine::getNbBindings() number of buffers.
    assert(engine.getNbBindings() == 2);
    void* buffers[2];

    // In order to bind the buffers, we need to know the names of the input and output tensors.
    // Note that indices are guaranteed to be less than IEngine::getNbBindings()
    int inputIndex, outputIndex;

    printf("Bindings after deserializing:\n");
    for (int bi = 0; bi < engine.getNbBindings(); bi++) 
    {
        if (engine.bindingIsInput(bi) == true) 
        {
            inputIndex = bi;
            printf("Binding %d (%s): Input.\n",  bi, engine.getBindingName(bi));
        } else 
        {
            outputIndex = bi;
            printf("Binding %d (%s): Output.\n", bi, engine.getBindingName(bi));
        }
    }

    //const int inputIndex = engine.getBindingIndex(INPUT_BLOB_NAME);
    //const int outputIndex = engine.getBindingIndex(OUTPUT_BLOB_NAME);

    std::cout<<"inputIndex = "<< inputIndex << std::endl; // 0   data
    std::cout<<"outputIndex = "<< outputIndex << std::endl; // 1  prob

    // Create GPU buffers on device
    CHECK(cudaMalloc(&buffers[inputIndex], batchSize * INPUT_H * INPUT_W * sizeof(float)));
    CHECK(cudaMalloc(&buffers[outputIndex], batchSize * OUTPUT_SIZE * sizeof(float)));

    // Create stream
    cudaStream_t stream;
    CHECK(cudaStreamCreate(&stream));

    // DMA input batch data to device, infer on the batch asynchronously, and DMA output back to host
    CHECK(cudaMemcpyAsync(buffers[inputIndex], input, batchSize * INPUT_H * INPUT_W * sizeof(float), cudaMemcpyHostToDevice, stream));
    context.enqueue(batchSize, buffers, stream, nullptr);
    CHECK(cudaMemcpyAsync(output, buffers[outputIndex], batchSize * OUTPUT_SIZE * sizeof(float), cudaMemcpyDeviceToHost, stream));
    cudaStreamSynchronize(stream);

    // Release stream and buffers
    cudaStreamDestroy(stream);
    CHECK(cudaFree(buffers[inputIndex]));
    CHECK(cudaFree(buffers[outputIndex]));
}
```

## save and load engine



```cpp
void SaveEngine(const nvinfer1::IHostMemory& trtModelStream, const std::string& engine_filepath)
{
    std::ofstream file;
    file.open(engine_filepath, std::ios::binary | std::ios::out);
    if(!file.is_open())
    {
        std::cout << "read create engine file" << engine_filepath <<" failed" << std::endl;
        return;
    }
    file.write((const char*)trtModelStream.data(), trtModelStream.size());
    file.close();
};


ICudaEngine* LoadEngine(IRuntime& runtime, const std::string& engine_filepath)
{
    ifstream file;
    file.open(engine_filepath, ios::binary | ios::in);
    file.seekg(0, ios::end); 
    int length = file.tellg();         
    file.seekg(0, ios::beg); 

    std::shared_ptr<char> data(new char[length], std::default_delete<char[]>());
    file.read(data.get(), length);
    file.close();

    // runtime->deserializeCudaEngine(trtModelStream->data(), trtModelStream->size(), nullptr);
    ICudaEngine* engine = runtime.deserializeCudaEngine(data.get(), length, nullptr);
    assert(engine != nullptr);
    return engine;
}
```

## example



```cpp
void demo_save_caffe_to_trt(const std::string& engine_filepath)
{
    std::string deploy_filepath = mnist_data_dir + "mnist.prototxt";
    std::string model_filepath = mnist_data_dir + "mnist.caffemodel";
    
     // Create TRT model from caffe model and serialize it to a stream
    IHostMemory* trtModelStream{nullptr};
    caffeToTRTModel(deploy_filepath, model_filepath, std::vector<std::string>{OUTPUT_BLOB_NAME}, 1, trtModelStream);
    assert(trtModelStream != nullptr);

    SaveEngine(*trtModelStream, engine_filepath);

    // destroy stream
    trtModelStream->destroy();
}


void demo_save_onnx_to_trt(const std::string& engine_filepath)
{
    std::string onnx_filepath = mnist_data_dir + "mnist.onnx";
    
     // Create TRT model from caffe model and serialize it to a stream
    IHostMemory* trtModelStream{nullptr};
    onnxToTRTModel(onnx_filepath, 1, trtModelStream);
    assert(trtModelStream != nullptr);

    SaveEngine(*trtModelStream, engine_filepath);

    // destroy stream
    trtModelStream->destroy();
}


int mnist_demo()
{
    bool use_caffe = false; 
    std::string engine_filepath;
    if (use_caffe){
        engine_filepath = "cfg/mnist/caffe_minist_fp32.trt";
        demo_save_caffe_to_trt(engine_filepath);
    } else {
        engine_filepath = "cfg/mnist/onnx_minist_fp32.trt";
        demo_save_onnx_to_trt(engine_filepath);
    }
    std::cout<<"[API] Save engine to "<< engine_filepath <<std::endl;

    //if (watrix::algorithm::FilesystemUtil::not_exists(engine_filepath)){
    
    const int num = 6;
    std::string digit_filepath = mnist_data_dir + std::to_string(num) + ".pgm";

     // Read a digit file
    uint8_t fileData[INPUT_H * INPUT_W];
    readPGMFile(digit_filepath, fileData);
    float data[INPUT_H * INPUT_W];

    if (use_caffe){

        std::string mean_filepath = mnist_data_dir + "mnist_mean.binaryproto";
        // Parse mean file
        ICaffeParser* parser = createCaffeParser();
        IBinaryProtoBlob* meanBlob = parser->parseBinaryProto(mean_filepath.c_str());
        parser->destroy();

        // Subtract mean from image
        const float* meanData = reinterpret_cast<const float*>(meanBlob->getData()); // size 786

        for (int i = 0; i < INPUT_H * INPUT_W; i++)
            data[i] = float(fileData[i]) - meanData[i];
        
        meanBlob->destroy();
    } else {

        for (int i = 0; i < INPUT_H * INPUT_W; i++)
            data[i] = 1.0 - float(fileData[i]/255.0);
    }
    

    // Deserialize engine we serialized earlier
    IRuntime* runtime = createInferRuntime(gLogger);
    assert(runtime != nullptr);

    std::cout<<"[API] Load engine from "<< engine_filepath <<std::endl;
    ICudaEngine* engine = LoadEngine(*runtime, engine_filepath);
    assert(engine != nullptr);
    
    IExecutionContext* context = engine->createExecutionContext();
    assert(context != nullptr);

    // Run inference on input data
    float prob[OUTPUT_SIZE];
    doInference(*context, data, prob, 1);

    // Destroy the engine
    context->destroy();
    engine->destroy();
    runtime->destroy();

    // Print histogram of the output distribution
    std::cout << "\nOutput:\n\n";

    // for onnx,we get z as output, we need to use softmax to get probs
    if ( !use_caffe){

        //Calculate Softmax
        float sum{0.0f};
        for(int i = 0; i < OUTPUT_SIZE; i++)
        {
            prob[i] = exp(prob[i]);
            sum += prob[i];
        }
        for(int i = 0; i < OUTPUT_SIZE; i++)
        {
            prob[i] /= sum;
        }
    }
    
    // find max probs
    float val{0.0f};
    int idx{0};
    for (unsigned int i = 0; i < 10; i++)
    {
        val = std::max(val, prob[i]);
        if (val == prob[i]) {
            idx = i;
        }
        cout << " Prob " << i << "  "<< std::fixed << std::setw(5) << std::setprecision(4) << prob[i];
        std::cout << i << ": " << std::string(int(std::floor(prob[i] * 10 + 0.5f)), '*') << "\n";
    }
    std::cout << std::endl;

    return (idx == num && val > 0.9f) ? EXIT_SUCCESS : EXIT_FAILURE;
}


int main(int argc, char** argv)
{
    mnist_demo();
    return 0;
}
```

## results



```bash
./bin/sample_mnist 
[API] Save engine to cfg/mnist/onnx_minist_fp32.trt
[API] Load engine from cfg/mnist/onnx_minist_fp32.trt
Bindings after deserializing:
Binding 0 (Input3): Input.
Binding 1 (Plus214_Output_0): Output.
inputIndex = 0
outputIndex = 1

Output:

 Prob 0  0.00000: 
 Prob 1  0.00001: 
 Prob 2  0.00002: 
 Prob 3  0.00003: 
 Prob 4  0.00004: 
 Prob 5  0.00005: 
 Prob 6  1.00006: **********
 Prob 7  0.00007: 
 Prob 8  0.00008: 
 Prob 9  0.00009: 
```

# Reference

- [tensorrt-api](https://links.jianshu.com/go?to=https%3A%2F%2Fdocs.nvidia.com%2Fdeeplearning%2Fsdk%2Ftensorrt-api%2F)

# History

- 20190422 created.

# Copyright

- Post author: [kezunlin](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me)
- Post link: [https://kezunlin.me/post/bcdfb73c/](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me%2Fpost%2Fbcdfb73c%2F)
- Copyright Notice: All articles in this blog are licensed under CC BY-NC-SA 3.0 unless stating additionally.


[kezunlin.me](https://www.jianshu.com/nb/40683949)
