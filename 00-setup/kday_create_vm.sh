#!/usr/bin/env bash
#Download "K-Day 01" VirtualBox image (vdi).
#Create a folder named "K-Day 01" and copy the K-Day 01.vdi to the folder.

#Create VirtualBox entry.
VBoxManage createvm --name "K-Day 01" --ostype Ubuntu_64 --register
#Create ssh port forwarding.
VBoxManage modifyvm "K-Day 01" --ostype Ubuntu_64 --cpus 2 --memory 4000 --natpf1 "guestssh,tcp,,2222,,22"
#Create the sata controller.
VBoxManage storagectl "K-Day 01"  --name hd1 --add sata --portcount 2
#Point the downloaded ubuntu vdi file using media switch.
VboxManage storageattach "K-Day 01" --storagectl hd1 --port 1 --type hdd --medium  "/Users/deep/VirtualBox VMs/K-Day 01/K-Day 01.vdi" --setuuid ""
#Start the vm
VBoxManage startvm "K-Day 01"
