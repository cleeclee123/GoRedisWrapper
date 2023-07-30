package main

import (
  "context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
  "github.com/joho/godotenv"
	"net/http"
  "fmt"
  "os"
)

func main() {
  client := redis.NewClient(&redis.Options{
    Addr: goDotEnv("REDIS_ADDY"),
    Password: goDotEnv("REDIS_PW"), 
    DB:       0,  
  })
  ctx := context.Background()

  err := client.Ping(ctx).Err()
  if err != nil {
    panic("redis client connect failed")
  }

	router := gin.Default()
  router.GET("/:apikey", home)
  router.GET("/set/:key/:value/:apikey", setRedis(ctx, client))
  router.GET("/get/:key/:apikey", getRedis(ctx, client))

	router.Run("localhost:8080")
}

func goDotEnv(key string) string {
  err := godotenv.Load(".env")
  if err != nil {
    panic("env load failed")
  }

  return os.Getenv(key)
}

func home(c *gin.Context) {
  apikey := c.Param("apikey")

  if (apikey != goDotEnv("API_KEY")) {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad api key"})
    return
  }
  c.IndentedJSON(http.StatusNotFound, gin.H{"hello": "welcome to my basic redis rest wrapper"})
}

func setRedis(ctx context.Context, client *redis.Client) gin.HandlerFunc {
  fn := func(c *gin.Context) {
    key := c.Param("key")
    value := c.Param("value")
    apikey := c.Param("apikey")

    if (apikey != goDotEnv("API_KEY")) {
      c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad api key"})
      return
    }
    
    err := client.Set(ctx, key, value, 0).Err()
    if err != nil {
      c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
      return
    }
    c.IndentedJSON(http.StatusOK, gin.H{"success": 200})
  }

  return gin.HandlerFunc(fn)
}

func getRedis(ctx context.Context, client *redis.Client) gin.HandlerFunc {
  fn := func(c *gin.Context) {
    key := c.Param("key")
    apikey := c.Param("apikey") 
    
    if (apikey != goDotEnv("API_KEY")) {
      c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad api key"})
      return
    }
    
    val, err := client.Get(ctx, key).Result()
    fmt.Println("value", val)
    if err != nil {
      c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
      return
    }
    c.IndentedJSON(http.StatusOK, gin.H{key: val})
  }
  
  return gin.HandlerFunc(fn)
}
