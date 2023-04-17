package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	openai "github.com/sashabaranov/go-openai"
)

type ChatCompletionRequest struct {
  Messages []openai.ChatCompletionMessage `json:"messages"`
}

func main() {
  client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
  app := fiber.New()

  app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
  }))

  app.Post("/api/plan", func(c *fiber.Ctx) error {
    body := new(ChatCompletionRequest)
    if err := c.BodyParser(body); err != nil {
      return c.JSON(fiber.Map{
        "message": "",
        "error": err.Error(),
      })
    }

    fmt.Printf("body: %v\n", body.Messages)

    messages := append([]openai.ChatCompletionMessage{
      {
          Role:    openai.ChatMessageRoleSystem,
          Content: "You are a helpful assistant that help user to create activity plan on the weekend. And response it as list.",
      },
    }, body.Messages...)

    fmt.Printf("%v\n", messages)

    req := openai.ChatCompletionRequest{
      Model: openai.GPT3Dot5Turbo,
      Messages: messages,
    }

    resp, err := client.CreateChatCompletion(context.Background(), req)

    if err != nil {
      c.JSON(fiber.Map{
        "message": "",
        "error": err.Error(),
      })
    }

    return c.JSON(fiber.Map{
      "message": resp.Choices[0].Message.Content,
      "error": "",
    })
  })

  app.Listen(":4000")
}

