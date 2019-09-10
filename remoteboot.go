package remoteboot

import (
        "fmt"
        "github.com/spf13/cobra"
        "opendev.org/airship/airshipctl/pkg/environment"
        "os"
        "bytes"
        "encoding/json"
        "net/http"
)

type f struct {
        OdataType         string `json:"@odata.type"`
        Name              string `json:"Name"`
        MembersOdataCount int    `json:"Members@odata.count"`
        Members           []struct {
                OdataID string `json:"@odata.id"`
        } `json:"Members"`
        OdataContext     string `json:"@odata.context"`
        OdataID          string `json:"@odata.id"`
        RedfishCopyright string `json:"@Redfish.Copyright"`
}

var actionstr string
var endpoint string
var transport string
var targethost string
var hostid string

var f2 = new(f)
func getJson(url string, target interface{}) error {
        var client http.Client
        r, err := client.Get(url)
        if err != nil {
                return err
        }
        defer r.Body.Close()

        return json.NewDecoder(r.Body).Decode(target)
}

func NewRemoteCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
        remotebootstrapRootCmd := &cobra.Command{
                Use:   "remoteboot",
                Short: "remoteboot airshipctl",
                Run: func(cmd *cobra.Command, args []string) {
                        name, _ := cmd.Flags().GetString("name")
                        if name == "" {
                                name = "remote"
                        }

                        getJson(transport+"://"+endpoint+"/redfish/v1/Systems", f2)
                        p := &hostid
                        *p = f2.Members[0].OdataID
                        fmt.Println(hostid)
                        fmt.Println("Hi, please use remoteboot <redfishapi/ansible> subcommands, specify --action flag")
                },
        }

        remotetype2 := &cobra.Command{
                Use:   "redfishapi",
                Short: "redfishapi",
                Run:   emptyRun1}
        remotetype3 := &cobra.Command{
                Use:   "ansible",
                Short: "ansible",
                Run:   emptyRun2}

        remotetype2.AddCommand(NewRemoteTarget(rootSettings))
        remotetype3.AddCommand(NewRemoteTarget(rootSettings))
        remotebootstrapRootCmd.AddCommand(remotetype2, remotetype3, NewRemoteTarget(rootSettings))
        remotebootstrapRootCmd.PersistentFlags().StringVar(&actionstr, "action", "On", "ForceOff or On")
        //RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rackhdcli.yaml)")
        remotetype2.PersistentFlags().StringVar(&endpoint, "endpoint", "10.118.135.27:8000", "API endoint of Redfishapi")
        remotetype2.PersistentFlags().StringVar(&transport, "transport", "http", "http or https")

        return remotebootstrapRootCmd
}

func emptyRun1(cmd *cobra.Command, args []string) {
        out := cmd.OutOrStdout()
        fmt.Fprintf(out, "Running redfishapi, please specify  --endpoint flag and target host using target subcommand\n")
}

func emptyRun2(cmd *cobra.Command, args []string) {
        out := cmd.OutOrStdout()
        fmt.Fprintf(out, "Running ansible, please specify target host using target subcommand\n")
}

func NewRemoteTarget(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
        remotetarget := &cobra.Command{
                Use:   "target",
                Short: "target",
                Run: func(cmd *cobra.Command, args []string) {
                        name, _ := cmd.Flags().GetString("target")
                        if name == "" {
                                name = "target"
                                fmt.Println("You did not specify the target, please specify host using -T option")
                                os.Exit(1)

                        }
                        fmt.Println("Target " + name + ", specified but not used for now")
                        targethost = name
                        TargetRun()
                },
        }

        remotetarget.Flags().StringP("target", "T", "", "Set Target of exection")
        return remotetarget
}

func TargetRun() {
        fmt.Println(endpoint, transport)
        fmt.Println(targethost, actionstr)
        fmt.Println("remoteboot running on " + hostid)

        type Payload struct {
                ResetType string `json:"ResetType"`
        }

        data := Payload{
                ResetType: actionstr,
        }
        payloadBytes, err := json.Marshal(data)
        if err != nil {
                // handle err
        }
        body := bytes.NewReader(payloadBytes)

        getJson(transport+"://"+endpoint+"/redfish/v1/Systems", f2)
        p := &hostid
        *p = f2.Members[0].OdataID
        fmt.Println(hostid)
        url1 := transport+"://"+endpoint+hostid+"/Actions/ComputerSystem.Reset/"
        fmt.Println(url1)

        req, err := http.NewRequest("POST", url1, body)
        if err != nil {
                // handle err
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
                // handle err
        }
        defer resp.Body.Close()

}
