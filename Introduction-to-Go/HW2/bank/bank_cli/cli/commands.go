package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

type Command struct {
	rootCmd *cobra.Command
}

func (cmd *Command) Init() {
	cmd.rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "CLI client for account management",
	}

	cmd.rootCmd.AddCommand(cmd.getCmt())
	cmd.rootCmd.AddCommand(cmd.createCmt())
	cmd.rootCmd.AddCommand(cmd.updateNameCmt())
	// cmd.rootCmd.AddCommand(cmd.updateBalanceCmt())
	cmd.rootCmd.AddCommand(cmd.deleteCmt())
}

func (cmd *Command) Execute() {
	cmd.rootCmd.Execute()
}

func (cmd *Command) getCmt() *cobra.Command {
	return &cobra.Command{
		Use:   "get [id]",
		Short: "Get account by id",
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			resp, err := http.Get("http://localhost:8080/account/" + id)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
		},
	}
}

func (cmd *Command) createCmt() *cobra.Command {
	return &cobra.Command{
		Use:   "create [id] [name] [balance]",
		Short: "Create a new account",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			name := args[1]
			balance, _ := strconv.ParseFloat(args[2], 64)

			type Account struct {
				ID      string  `json:"id"`
				Name    string  `json:"name"`
				Balance float64 `json:"balance"`
			}

			account := Account{
				ID: id,
				Name: name,
				Balance: balance,
			}
			body, _ := json.Marshal(account)
			resp, err := http.Post("http://localhost:8080/account", "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()
			respBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(respBody))
		},
	}
}

func (cmd *Command) updateNameCmt() *cobra.Command {
	return &cobra.Command{
		Use:   "update-name [id] [name]",
		Short: "Update account name",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			name := args[1]

			reqBody := map[string]string{"name": name}
			body, _ := json.Marshal(reqBody)
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/account/"+id+"/name", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			defer resp.Body.Close()
			respBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(respBody))
		},
	}
}

func (cmd *Command) deleteCmt() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete account by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			client := &http.Client{}
			req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/account/"+id, nil)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNoContent {
				fmt.Println("Account deleted successfully")
			} else {
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println(string(body))
			}
		},
	}
}
