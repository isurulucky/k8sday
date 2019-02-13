#!/usr/bin/env bash
#Create folder
mkdir ~/kday
#Download "K-Day 01" VirtualBox image (vdi).
wget http://atuwa.private.wso2.com/VMs/K-Day%2001.vdi.tgz -P ~/kday
cd ~/kday
tar -xzvf K-Day%2001.vdi.tgz
#Create VirtualBox entry.
VBoxManage createvm --name "K-Day 01" --ostype Ubuntu_64 --register
#Create ssh port forwarding.
VBoxManage modifyvm "K-Day 01" --ostype Ubuntu_64 --cpus 2 --memory 4000 --natpf1 "guestssh,tcp,,2222,,22"
#Create the sata controller.
VBoxManage storagectl "K-Day 01"  --name hd1 --add sata --portcount 2
#Point the downloaded ubuntu vdi file using media switch.
VboxManage storageattach "K-Day 01" --storagectl hd1 --port 1 --type hdd --medium  "~/kday/K-Day 01.vdi" --setuuid ""
#Start the vm
VBoxManage startvm "K-Day 01"
