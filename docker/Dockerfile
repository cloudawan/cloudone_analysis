FROM ubuntu:14.04

# Install curl
RUN apt-get install -y curl

# Bundle app source
COPY . /src

# Copy configuration
RUN mkdir -p /etc/cloudone_analysis
RUN mkdir -p /var/log/cloudone_analysis
RUN cp /src/cloudone_analysis/configuration/* /etc/cloudone_analysis/

# Privilege
RUN chmod a+x /src/run.sh
RUN chmod a+x /src/cloudone_analysis/cloudone_analysis

EXPOSE  8082
CMD ["/src/run.sh"]
