package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

var adminToken = os.Getenv("ADMINTOKEN")

func parseIP(c *fiber.Ctx, header string) net.IP {
	var ip string
	var requestedIP = c.Get(header, "")
	if (requestedIP == "" || requestedIP == "mine") {
		var ips []string = c.IPs() // X-Forwarded-For

		if (len(ips) > 0) {
			ip = ips[len(ips) - 1] // last item since this is sitting behind heroku
		} else {
			ip = c.IP()
		}
	} else {
		ip = requestedIP
	}
	return net.ParseIP(utils.ImmutableString(ip))
}

func parsePort(c *fiber.Ctx, header string) int {
	var requestedPort = c.Get(header, "")
	if requested, err := strconv.Atoi(requestedPort); err == nil {
		return requested
	} else if (requestedPort == "") {
		return 6969
	}
	return 0
}

func main() {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		var authorization = c.Get("Authorization", "")
		if (authorization == "") {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Missing authorization token",
			})
		} else if (authorization != "Bearer " + adminToken) {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token",
			})
		}

		var inputIP = parseIP(c, "X-Input-Ip")
		var outputIP = parseIP(c, "X-Output-Ip")
		if (inputIP == nil || outputIP == nil) {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid IP",
			})
		}

		var inputPort = parsePort(c, "X-Input-Port")
		var outputPort = parsePort(c, "X-Output-Port")
		if (inputPort < 1000 || inputPort > 65535 ||
				outputPort < 1000 || outputPort > 65535) {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Invalid port: it has to be in range [1000-65535]",
			})
		}

		c.Status(200).JSON(fiber.Map{
			"success": true,
			"input": fiber.Map{
				"ip": inputIP.String(),
				"port": inputPort,
			},
			"output": fiber.Map{
				"ip": outputIP.String(),
				"port": outputPort,
			},
		})

		go StartConnection(Address{
			ip: inputIP.String(),
			port: inputPort,
		}, Address{
			ip: outputIP.String(),
			port: outputPort,
		})
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
