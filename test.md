* signin，记得保存 token

```shell
http://localhost:8080/simpleblog/user/signin?username=user5&password=pass5
```

* getArticleById

```
http://localhost:8080/simpleblog/user/article/1
```

* getArticles

```
http://localhost:8080/simpleblog/user/articles?page=1
```

* deleteArticle

```
http://localhost:8080/simpleblog/user/deleteArticle/1
```

* createComment 
  * 将之前登陆的 user 的 token 放在 Authorization 后，author 对应登陆的 user

```
curl -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU2OTE4MzcsImlhdCI6MTU3NTY4ODIzN30.4Uw5rYZPCxB7uXNrVtn69Tmsy-831CgPnF8e555z-ko" http://localhost:8080/simpleblog/user/article/2/comment -X POST -d '{"content":"new content3","author":"user5"}'
```

* getComments

```
http://localhost:8080/simpleblog/user/article/2/comments
```



