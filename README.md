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


