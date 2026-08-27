package main

import (
	"flag"
	"fmt"
	"os"

	"campgoods/catalog"
	"campgoods/query"
	"campgoods/workflow"
)

func main() {
	path := flag.String("db", "campgoods.db", "path to the camp goods database")
	flag.Parse()
	app, err := workflow.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer app.Close()
	if flag.NArg() == 0 {
		fmt.Println("campgoods ready: use register, receive, price, or list")
		return
	}
	if err := runCommand(app, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runCommand(app *workflow.App, args []string) error {
	switch args[0] {
	case "register":
		if len(args) != 8 {
			return fmt.Errorf("register expects id sku name category unit price description")
		}
		var price int64
		if _, err := fmt.Sscan(args[6], &price); err != nil {
			return err
		}
		product, err := app.RegisterProduct(args[1], catalog.ProductDraft{SKU: args[2], Name: args[3], Category: args[4], Unit: args[5], PriceCents: price, Description: args[7]})
		if err != nil {
			return err
		}
		fmt.Printf("registered %s %s\n", product.ID, product.Name)
		return nil
	case "receive":
		if len(args) != 5 {
			return fmt.Errorf("receive expects movement product quantity reason")
		}
		var quantity int64
		if _, err := fmt.Sscan(args[3], &quantity); err != nil {
			return err
		}
		result, err := app.ReceiveStock(args[1], args[2], quantity, args[4])
		if err == nil {
			fmt.Println(result)
		}
		return err
	case "price":
		if len(args) != 5 {
			return fmt.Errorf("price expects change product cents reason")
		}
		var cents int64
		if _, err := fmt.Sscan(args[3], &cents); err != nil {
			return err
		}
		result, err := app.ChangePrice(args[1], args[2], cents, args[4])
		if err == nil {
			fmt.Println(result)
		}
		return err
	case "list":
		page, err := query.ParsePage(valueOrDefault(args, 1))
		if err != nil {
			return err
		}
		view, err := app.ListProducts(query.ListRequest{Page: page, PageSize: 10})
		if err == nil {
			fmt.Println(view.Render())
		}
		return err
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func valueOrDefault(values []string, index int) string {
	if index < len(values) {
		return values[index]
	}
	return ""
}
