# Config IAA Workqueue
---

## Install idxd-config
First of all, you need root privileges.

```shell
sudo su
```

Follow [https://github.com/intel/idxd-config](https://github.com/intel/idxd-config) to install idxd-config .

## Check workqueues

```shell
accel-config list
```
If you see something like:
```json
[
  {
    "dev":"dsa0",
    "read_buffer_limit":0,
    "max_groups":4,
    "max_work_queues":8,
    "max_engines":4,
    "work_queue_size":128,
    "numa_node":0,
    ...
  }
]
```
Your system administrator have already configed workqueues for you, you can skip this guide.

## Check your devices

```shell
ls /sys/bus/dsa/devices
```

If you have devices on the machine, you can see some directories like: `iax1` or  `dsa0`, those directories represent your IAA/DSA devices.

## Config a workqueue

Config `iax1` workqueue `wq1.0` :

| field          | value  |
| :------------- | ------ |
| mode           | shared |
| group          | 0      |
| size           | 16     |
| priority       | 10     |
| block on fault | true   |
| type           | user   |
| name           | app1   |
| driver         | user   |

```shell

accel-config config-wq iax1/wq1.0 -g 0 -t 15 -s 16 -p 10 -b 1 -m shared -y user -n app1 -d user

accel-config config-engine iax1/engine1.0 -g 0

accel-config enable-device iax1

accel-config enable-wq iax1/wq1.0

```

You can use `accel-config help config-wq` to learn how to use `accel-config`.

## Test if it works

Save the code snippet as `main.go`

```go
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/intel/ixl-go/crc"
)

func main() {
	data := make([]byte, 128)
	rand.Read(data)
	if crc.Ready() {
		c, _ := crc.NewCalculator()
		value, err := c.CheckSum64(data, crc.ISO)
		if err != nil {
			fmt.Println("occurs error:", err)
		} else {
			fmt.Println(hex.EncodeToString(data))
			fmt.Println("crc64 ISO:", value)
		}

	} else {
		fmt.Println("No IAA devices found")
	}
}

```

Run It!

```
# go run main.go
94a79a59365343212a0d51af48e078bba88485e38e4fc96a4d375343a81b93f110bc9df52c02f98e3abf6b63e6dc0cb3923afb6f12ccf5141cae615f846128e5585588d3bb88bcb493f9a0f6aa45077943da21dcfe52911c70e3b04c2c7423acf2f1c1783ec1c595682a9152c902a410af5407620b323e407cb5d0fc30e96911
crc64 ISO: 10562025956998192512
```