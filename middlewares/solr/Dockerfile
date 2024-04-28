# docker build -t atcodersearch-solr:latest .
# docker container run --rm --name solr -p 8983:8983 atcodersearch-solr:latest

FROM solr:9.1.0

USER root
RUN sed -i -E 's/^#(ja_JP.UTF-8)/\1/' /etc/locale.gen \
    && locale-gen
ENV LANG ja_JP.UTF-8

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV SOLR_TIMEZONE = "JST"

COPY --chown=solr:solr ./problem /var/solr/data/problem
COPY --chown=solr:solr ./user /var/solr/data/user
COPY --chown=solr:solr ./submission /var/solr/data/submission
COPY --chown=solr:solr ./solution /var/solr/data/solution
COPY --chown=solr:solr ./dict/lucene-analysis-kuromoji-9.3.0-unidic-2.1.2.jar /opt/solr/server/solr-webapp/webapp/WEB-INF/lib/lucene-analysis-kuromoji-9.3.0.jar

USER solr
