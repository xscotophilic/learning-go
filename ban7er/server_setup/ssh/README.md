# Installing and Configuring OpenSSH Server in VirtualBox

## Step 1: Install OpenSSH Server

To install the OpenSSH server on an Ubuntu system, run the following commands:

```shell
sudo apt update
sudo apt install openssh-server -y
```

## Step 2: Enable and Start the SSH Service

Once installed, enable and start the SSH service using the following commands:

```shell
sudo systemctl enable ssh
sudo systemctl start ssh
```

## Step 3: Verify SSH Service Status

To check whether the SSH service is running, execute:

```shell
sudo systemctl status ssh
```

## Step 4: Configure Network Settings in VirtualBox

To allow SSH access from the host system to the VirtualBox guest, configure port forwarding as follows:

1. Open VirtualBox.
2. Select your VM and go to **Settings** > **Network**.
3. Select **Adapter** (Ensure it is attached to **NAT**).
4. Click **Advanced** > **Port Forwarding**.
5. Add a new rule with the following settings:
   - **Protocol:** TCP
   - **Host IP:** `127.0.0.1`
   - **Host Port:** `2222`
   - **Guest IP:** Leave blank
   - **Guest Port:** `22`
6. Click **OK** to save the settings.

## Step 5: Connect to the Virtual Machine via SSH

After configuring port forwarding, you can SSH into the guest VM from the host machine using:

```shell
ssh -p 2222 username@127.0.0.1
```

Replace `username` with your actual VM username.
