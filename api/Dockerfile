ARG DOCKER_BASE_IMAGE

FROM $DOCKER_BASE_IMAGE

ARG IMAGE_NAME
ARG NAME

COPY datamodel_installer.sh /datamodel_installer.sh
COPY build/crds /crds

RUN echo $IMAGE_NAME > /IMAGE
RUN echo $NAME > /NAME

CMD /datamodel_installer.sh
