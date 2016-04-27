package legacy

import (
	"fmt"
	//	"os"
	"os/exec"
	"trustcloud/util"
)

func renderCard(cardDetails *CardDetails) int {
	util.InfoLog.Println("Render card ", cardDetails)

	//out, err := exec.Command("sh", "-c", "/bin/services/tst.sh", os.Getenv("PATH")).Output()
	// Working
	//	cmd := exec.Cmd{
	//		Path: "CutyCapt",
	//		Args: []string{"tst.sh"},
	//		Dir:  "services/cutycapt",
	//	}
	//
	//	cmd := exec.Cmd{
	//		Path: "xvfb-run",
	//		Args: []string{"-a", "--server-args=\"-screen 0, 640x480x24\"", "services/cutycapt/CutyCapt"},
	//		Dir:  "services/cutycapt",
	//	}

	//xvfb-run -a --server-args="-screen 0, 640x480x24" ' . FS_ROOT_INCLUDE . 'lib/cutycapt/CutyCapt

	/*
		xvfb-run -a --server-args="-screen 0, 640x480x24" /var/www/api.trustcloud.com/include/lib/cutycapt/CutyCapt --url=https://api.trustcloud.com/display/renderidcard?svg=no\&allowUpdate=no\&userid=trustcloud-180035\&size=b3 --out=/tmp/1o.png --zoom-factor=2.0 --min-width=100 --min-height=100
	*/

	//	cmd := exec.Command("sh",
	//		"-c", "services/tst.sh")

	cmd := exec.Command("xvfb-run",
		"-a", "services/cutycapt/tst.sh")

	//	cmd := exec.Command("xvfb-run",
	//		"-a", "services/cutycapt/CutyCapt") // --url=http://www.example.org/ --out=localfile.png")
	//		"--server-args=\"-screen 0, 640x480x24\"", "services/cutycapt/CutyCapt", "--url=https://api.trustcloud.com/display/renderidcard?svg=no&allowUpdate=no&userid=trustcloud-180035&size=b3 --out=/tmp/1o.png --zoom-factor=2.0 --min-width=100 --min-height=100")
	//	cmd.Dir = "services"

	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(out))
		return 0
	} else {
		fmt.Println(string(out))
	}

	//out, err := exec.Command("./services/cutycapt/CutyCapt").Output()

	//	util.InfoLog.Println("out... ", string(out))
	//
	//	util.InfoLog.Println("err... ", err)

	return 0
}
