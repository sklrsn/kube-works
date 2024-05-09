FROM rockylinux:9.3 AS base 

RUN yum update -y && yum install -y httpd

COPY index.html /var/www/html/
COPY ./httpd.conf /usr/local/apache2/conf/httpd.conf

EXPOSE 80

CMD ["/usr/sbin/httpd", "-D", "FOREGROUND"]