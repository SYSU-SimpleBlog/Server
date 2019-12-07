/*
 * simple blog
 *
 * A Simple Blog
 *
 * API version: 1.0.0
 * Contact: apiteam@swagger.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package t

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//

	sw "github.com/SYSU-SimpleBlog/Server/go"

	"github.com/boltdb/bolt"
)

func CreateComment() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create
	err = db.Update(func(tx *bolt.Tx) error {
		a := tx.Bucket([]byte("Article"))
		c := a.Cursor()
		b := tx.Bucket([]byte("Comment"))

		if b == nil {
			b, err = tx.CreateBucket([]byte("Comment"))
			if err != nil {
				log.Fatal(err)
			}
		}
		if b != nil {
			var article sw.Article
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var comment sw.Comment
				err := json.Unmarshal(v, &article)
				if err != nil {
					return err
				}
				fmt.Println(article.Id)
				var id int
				id = int(article.Id)
				filePath := "./data/" + strconv.Itoa(id) + "/comments"
				files, err := ioutil.ReadDir(filePath)
				if err != nil {
					log.Fatal(err)
				}

				for i := 1; i <= len(files); i++ {
					file, err := os.OpenFile(filePath+"/"+files[i-1].Name(), os.O_RDWR, 0666)
					buf := bufio.NewReader(file)
					user, err := buf.ReadString('\n')
					user = strings.TrimSpace(user)
					fmt.Println(user)

					time, err := buf.ReadString('\n')
					time = strings.TrimSpace(time)
					fmt.Println(time)

					var content string
					for {
						line, err := buf.ReadString('\n')
						line = strings.TrimSpace(line)
						content = content + line
						if err != nil {
							if err == io.EOF {
								fmt.Println("File read ok!")
								break
							} else {
								fmt.Println("Read file error!", err)
							}
						}
					}

					//timeStr := time.Now().Format("2006-01-02 15:04:05")
					comment = sw.Comment{time, content, user, article.Id}
					fmt.Println(comment)
					vc, err := json.Marshal(comment)
					err = b.Put([]byte(strconv.Itoa(int(article.Id))+"_"+strconv.Itoa(i)), []byte(vc))
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		} else {
			return errors.New("Table Comment doesn't exist")
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func DBTestComment() {
	CreateComment()
}
