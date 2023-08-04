---
sidebar_position: 1
slug: /
---

# Introduction

The `migrations` package is an abstraction that can migrate anything. What you need is to implement a [Source](https://pkg.go.dev/github.com/jamillosantos/migrations#Source)
and a [Target](https://pkg.go.dev/github.com/jamillosantos/migrations#Target).

For now, there are a couple of implementations of Source and Target that you can use:
* [SQL](/sql/getting-started)
* [MongoDB](/mongo)
* [Functions](/functions)

