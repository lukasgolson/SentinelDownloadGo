/*
Copyright © 2024 Lukas G. Olson <olson@student.ubc.ca>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
)

// signupCmd represents the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Opens your web browser to the Sentinel Hub signup page.",

	Run: func(cmd *cobra.Command, args []string) {

		const url string = "https://documentation.dataspace.copernicus.eu/Registration.html"

		QR, _ := cmd.Flags().GetBool("QR")

		if QR {
			config := qrterminal.Config{
				HalfBlocks: true,
				Level:      qrterminal.L,
				Writer:     os.Stdout,
			}
			qrterminal.GenerateWithConfig(url, config)
			return
		} else {
			fmt.Println("Opening instructions in your default web-browser...")
			var err error
			switch runtime.GOOS {
			case "linux":
				err = exec.Command("xdg-open", url).Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
			case "darwin":
				err = exec.Command("open", url).Start()
			default:
				err = fmt.Errorf("unsupported platform")
			}
			if err != nil {
				fmt.Println(err)
			}
		}

	},
}

func init() {
	authCmd.AddCommand(signupCmd)

	signupCmd.Flags().BoolP("QR", "q", false, "Display a QR code of the signup page")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
