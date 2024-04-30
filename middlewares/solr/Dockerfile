# docker build -t atcodersearch-solr:latest .
# docker container run --rm --name solr -p 8983:8983 atcodersearch-solr:latest

FROM solr:9.5

USER root
RUN sed -i -E 's/^#(ja_JP.UTF-8)/\1/' /etc/locale.gen \
    && locale-gen
ENV LANG ja_JP.UTF-8

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV SOLR_TIMEZONE = "JST"

COPY --chown=solr:solr ./problems /var/solr/data/problems
COPY --chown=solr:solr ./users /var/solr/data/users
COPY --chown=solr:solr ./submissions /var/solr/data/submissions
COPY --chown=solr:solr ./solutions /var/solr/data/solutions
COPY --chown=solr:solr ./dict/lucene-analysis-kuromoji-9.9.2-unidic-2.1.2.jar /opt/solr/server/solr-webapp/webapp/WEB-INF/lib/lucene-analysis-kuromoji-9.9.2.jar

USER solr
