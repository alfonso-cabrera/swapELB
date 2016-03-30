# swapELB

Get registered EC2 instances from one ELB and programatically add them to another ELB.

Edit the main.go file to add the source and destination ELB names (not endpoints!) along with the IAM role that has elasticloadbalancing permissions.
