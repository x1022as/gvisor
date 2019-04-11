FROM ubuntu:16.04
#FROM ubuntu:16.04
ADD apt.conf /etc/apt/apt.conf
# Install sysbench
#ADD iozone.test /mnt/mvm-test/iozone.test
RUN apt-get update
RUN apt install --yes sysbench
#ADD testscript /home/testscript
ADD iozone /usr/bin/iozone
RUN rm -fv /etc/apt/apt.conf
CMD ["/bin/bash"]
