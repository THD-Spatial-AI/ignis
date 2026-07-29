[1mdiff --git a/environment/Dockerfile b/environment/Dockerfile[m
[1mindex d192974..1d74b86 100644[m
[1m--- a/environment/Dockerfile[m
[1m+++ b/environment/Dockerfile[m
[36m@@ -3,7 +3,8 @@[m
 # ---------------------------------------------------------------------------[m
 # Stage 1: builder — has the full Go toolchain + source, never shipped[m
 # ---------------------------------------------------------------------------[m
[31m-FROM golang:1.26.5-bookworm AS builder[m
[32m+[m[32mARG GO_VERSION=1.26.5[m
[32m+[m[32mFROM golang:${GO_VERSION}-bookworm AS builder[m
 [m
 WORKDIR /app[m
 [m
[1mdiff --git a/environment/Dockerfile.build_db b/environment/Dockerfile.build_db[m
[1mindex e2e935c..6b119d5 100644[m
[1m--- a/environment/Dockerfile.build_db[m
[1m+++ b/environment/Dockerfile.build_db[m
[36m@@ -9,7 +9,8 @@[m
 # ---------------------------------------------------------------------------[m
 # Stage 1: builder[m
 # ---------------------------------------------------------------------------[m
[31m-FROM golang:1.26.0-bookworm AS builder[m
[32m+[m[32mARG GO_VERSION=1.26.5[m
[32m+[m[32mFROM golang:${GO_VERSION}-bookworm AS builder[m
 [m
 WORKDIR /app[m
 COPY go.mod go.sum ./[m
