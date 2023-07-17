# docker build -t solr:atcoder-search .
# docker container run --rm --name solr -p 8983:8983 solr:atcoder-search

FROM solr:9.1.0

USER root
RUN sed -i -E 's/^#(ja_JP.UTF-8)/\1/' /etc/locale.gen \
    && locale-gen
ENV LANG ja_JP.UTF-8

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENV SOLR_TIMEZONE = "JST"

COPY --chown=solr:solr ./problem /var/solr/data/problem
COPY --chown=solr:solr ./recommend /var/solr/data/recommend
COPY --chown=solr:solr ./user /var/solr/data/user
COPY --chown=solr:solr ./dict/lucene-analysis-kuromoji-9.3.0-unidic-2.1.2.jar /opt/solr/server/solr-webapp/webapp/WEB-INF/lib/lucene-analysis-kuromoji-9.3.0.jar

USER solr
