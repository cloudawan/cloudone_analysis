FROM    centos:centos6

#EPEL repository
RUN rpm -iUvh http://dl.fedoraproject.org/pub/epel/6/x86_64/epel-release-6-8.noarch.rpm

# Bundle app source
COPY . /src

# Copy configuration
RUN mkdir -p /etc/cloudone_analysis
RUN mkdir -p /var/log/cloudone_analysis
RUN cp /src/cloudone_analysis/configuration/* /etc/cloudone_analysis/

# Priviledge
RUN chmod a+x /src/run.sh
RUN chmod a+x /src/cloudone_analysis/cloudone_analysis

EXPOSE  8082
CMD ["/src/run.sh"]
