sudo iptables -F
sudo iptables -A INPUT -p tcp --dport 14001 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport 14001 -j ACCEPT

sudo iptables -A INPUT -p tcp --dport 14002 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport 14002 -j ACCEPT

sudo iptables -A INPUT -p tcp --dport 14003 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport 14003 -j ACCEPT


sudo iptables -A INPUT -j DROP
echo "Plugged out network cable. Run 'sudo iptables -F' to flush the iptables"

