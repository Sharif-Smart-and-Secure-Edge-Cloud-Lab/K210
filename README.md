# MaixDuino

In this repository some guide into setting MaixDuino M1 board is provided. Maixduino consists of dual-core RISC-V processor and K210 AI Core. It is indeed a flexible and low power embedded system that can be used to develop AIoT applications.

## Setting Up
To start using MaixDuino you need a Type-C cable to connect it to your computer. Also a 24-pin LCD comes with the development board. Connect it to the main board. When you plug in the board, the LCD will display the following message:

![setting up](./imgs/setting-up.jpg)

This board comes with Micro Python installed. You can test it by serial port and running you scripts on it. 

First you need a serial terminal. In linux you can use `minicom` or `screen` to open a terminal. I used `screen` to open a terminal. Then you must find the serial port that is connected to your PC. Run the following command to find connected serial ports:

```
$ ls /dev/tty*
/dev/ttyUSB0  /dev/ttyUSB1
```

Then connect to the serial port using `screen`:

```
$ sudo screen /dev/ttyUSB0 115200 
```

Last argument shows the baud rate. Finally you will get a Python prompt. For more information refer to [MaixPy](https://wiki.sipeed.com/soft/maixpy/en/index.html) documentation.

## C++ baremetal programming
To run C++ baremetal programs first you need to download the SDKs. There are two SDKs provided to use with Kendryte K210: standalone and FreeRTOS. It seems that FreeRTOS SDK is deprecated so I didn't check it. Now I focuse on standalone SDK. Standalone SDK is in [this repository](https://github.com/kendryte/kendryte-standalone-sdk). Standalone SDK provides some libraries and drivers that you can use for building your programs.

First download the RISC-V toolchain that we will use it to compile the source codes. To do this head over to [this](https://github.com/kendryte/kendryte-gnu-toolchain) repository which is RISC-V toolchain for Kendryte devices. The installation procedure is pretty straight forward. However, I summarize the steps here. First clone the repository:

```bash
$ git clone --recursive https://github.com/kendryte/kendryte-gnu-toolchain
```
Before building the toolcahin you must install some prerequisites.

```bash
$ sudo apt-get install autoconf automake autotools-dev curl libmpc-dev libmpfr-dev libgmp-dev gawk build-essential bison flex texinfo gperf libtool patchutils bc zlib1g-dev libexpat-dev
```
Then run the following command to build the project:
    
```bash
$ ./configure --prefix=/opt/kendryte-toolchain --with-cmodel=medany --with-arch=rv64imafc --with-abi=lp64f
$ make -j8
```
`--prefix` argument shows where the toolchain will be installed. 

Now it's time to write a test code and compile it. After you downloaded the standalone SDK, you can put your codes in `src` folder and compile it. By default, there is a `hello_world` project which you can compile it to test the board and SDK. Compile it by running the following command:

```bash
$ mkdir build && cd build
$ cmake .. -DPROJ=<ProjectName> -DTOOLCHAIN=/opt/riscv-toolchain/bin && make
```
`-DTOOLCHAIN` shows the path to the toolchain which we installed it before. 

After compiling the project, you would find these two files: `hello_world` and `hello_world.bin`. `hello_world.bin` is the binary file that you can burn it into the flash memory of the board. `hello_world` is the `elf` file. Notice that you can't run it on your machine unless you have a RISC-V compatible machine. However, you can use `qemu` to run it on your machine.

Next, lets burn the binary to flash. Kendryte has developed a Python-based tool that uses UART to transfer `.bin` files to the board. I recommend you to make a virtual environment and install dependencies in it. Make a virtual environment by running the following command:

```bash
$ python3 -m venv <path to enviroment>
```
Then activate it by running:

```bash
$ source <path to enviroment>/bin/activate
```

You can either install all dependecies manually or use the `requirements.txt` file to install them. 

1. **Manual:**
    
    You can visit [this](https://github.com/kendryte/kflash.py) repository to find dependencies and kflash itself. Basically, you need to install `pyserial` and `pyelftools` to use kflash. Also install `kflash` itself from pypi. 
    ```bash
    pip install pyserial
    pip install pyelftools
    pip install kflash
    ```
2. **Requirements.txt:**
    
    Run the following command to install kflash and its dependencies:
    ```bash
    $ pip install -r requirements.txt
    ```
**Note:** Make sure that you have activated the virtual environment before running the previous commands.

Add your user to dialout group otherwise you would need to use `sudo` to run kflash:

```bash
$ sudo usermod -a -G dialout $(whoami)
```

Finally, transfer the `hello_world.bin` file to the board and run it:

```bash
$ kflash -B goE -p /dev/ttyUSB0 -t hello_world.bin
```

## Go baremetal programming
To use Go codes on Maixduio first you must install Go compiler on your machine. In order to build go from source, you need to have the go in first place! From version >= 1.4 this is necessary. If you don't have Go compiler already on you machine, then download latest relese from [go official website](https://go.dev/doc/install). After downloding the tarball, extract it and copy it to the `/usr/local/go` directory.

```bash
$ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz
```
Then add the go compiler to the PATH environment variable:

```bash
$ export PATH=$PATH:/usr/local/go/bin
```
Finally check if you have currectly configured the PATH variable:

```bash
$ go version
```

Now it is time to compile go from source to use it cross compiling. All steps are described in [this](https://embeddedgo.github.io/getting_started) page. First download the patch:

```bash
$ git clone https://github.com/embeddedgo/patch
```
Then download go compiler source code:
    
```bash
$ git clone https://go.googlesource.com/go goroot
```
Apply the patch:
    
```bash
$ cd goroot
$ git checkout go1.18.3
$ patch -p1 <../patch/go1.18.3
$ cd src
$ ./make.bash
```
You can use `./all.bash` instead of `./make.bash` to run all the tests after building the compiler. This takes extra time to build the compiler.

Now it's time to run a test code. I've written a code in `maix_blinky` folder. It is deriven from [this tutorial](https://embeddedgo.github.io/2020/05/31/bare_metal_programming_risc-v_in_go.html). First go to the `maix_blinky` directory. Then compile the code:

```bash
$ cd maix_blinky
$ go mod init maix_blinky
```
After that run the following command:

```bash
$ GOOS=noos GOARCH=riscv64 go build -tags k210 -ldflags '-M 0x80000000:6M'
```
You might get an error message about leds module. As the compiler says, install the module using the following command:

```bash
$ go get github.com/embeddedgo/kendryte/devboard/maixbit/board/leds
```
Then rerun the previous command. You would finally end up with these files:

```bash
$ ls 
go.mod  go.sum  main.go  maix_blinky
```
The `maix_blinky` is the ELF file. But you can't burn it into flash. So, you need to convert it to binary. Since it is compiled for RISC-V architecture, use the RISC-V toolchain to convert it to binary (refer to C++ baremetal programming for more details):

```bash
$ riscv64-unknown-elf-objcopy -O binary maix_blinky maix_blinky.bin
```
Now burn the code into the flash memory using Kflash:

```bash
$ kflash -B goE -p /dev/ttyUSB0 -t maix_blinky.bin
```
You can see that led is blinking! 

![blinky led](/imgs/go-blinky-led.gif)

## Developement with Arduino core
This section is primarily based on ([this page](https://maixduino.sipeed.com/en/get_started/install.html) .
To develop code with Arduino IDE, first install the Arduino IDE on your PC. Head over to [official Arduino download page](https://www.arduino.cc/en/Main/Software). Download the latest version and install it. To install on Linux, after untarring the tarball, run the following command:

```bash
$ sudo ./install.sh
```
After installing, you need to add your user to `dialout` group to grant access to serial ports. We've done it before (in C++ baremetal programming) but for sake of completeness, I'll mention it again:

```bash
$ sudo usermod -a -G dialout $(whoami)
```
Run the arduino IDE and select `File` -> `Preferences`. Add one of the following links to the `Additional Boards Manager URLs` section of the preferences:

```bash
http://dl.sipeed.com/MAIX/Maixduino/package_Maixduino_k210_index.json

## in case of slow download, try this:
http://dl.sipeed.com/MAIX/Maixduino/package_Maixduino_k210_dl_cdn_index.json
```

![arduino URLs](/imgs/arduino-URLs.png)

Now go to `Tools` -> `Board` -> `Boards manager` and search for `Maixduino`. Selecte the latest version and install it. 

After installing it's time to change board settings. First run a terminal. Make sure that you have activated the python virtual environment we built it before. Arduino IDE will use `kflash` to programm the board. Within that terminal, run the Arduino IDE (if you have installed `kflash` globally, you probably don't need to do that). Next in `Tools` menu, change these settings:

- `Board`: Choose your dev board (in our case choose `Sipeed Maixduino board` )
- `Burn Tool Frimware`: Choose `default`
- `Burn Baudrate`: Decrese the baudrate if download fails (I choosed 400MHz)
- `Port`: Serial port that the board is connected (e.g. /dev/ttyUSB0)
- `Programmer`: Burn tool. You **must** choose `kflash`

Finally, we are ready to programm the board using Arduino IDE! I have provided a test code in `arduino-test` directory. It's a simple blinker that you can see the result in the following video:

![blinker](./imgs/arduino-blinky-led.gif)

## Developement with PlatformIO