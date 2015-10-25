FROM    centos:centos6

#EPEL repository
RUN rpm -iUvh http://dl.fedoraproject.org/pub/epel/6/x86_64/epel-release-6-8.noarch.rpm

# Bundle app source
COPY . /src

# Copy configuration
RUN mkdir -p /etc/kubernetes_management_analysis
RUN mkdir -p /var/log/kubernetes_management_analysis
RUN cp /src/kubernetes_management_analysis/configuration/* /etc/kubernetes_management_analysis/

# Priviledge
RUN chmod a+x /src/run.sh
RUN chmod a+x /src/kubernetes_management_analysis/kubernetes_management_analysis

EXPOSE  8082
CMD ["/src/run.sh"]
