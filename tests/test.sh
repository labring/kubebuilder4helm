#!/bin/bash
# shellcheck disable=SC2164
cd helm-project
rm -rf *
kubebuilder4helm init
kubebuilder4helm create api --group user --version v1beta1 --kind Setting  --force
kubebuilder4helm create webhook --version v1beta1  --kind Setting --group user   --conversion --defaulting --programmatic-validation  --force
