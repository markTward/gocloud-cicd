service:
  gocloudAPI:
    image:
      repository: {{.Repo}}
      tag: {{.Tag}}
  gocloudGrpc:
    image:
      repository: {{.Repo}}
      tag: {{.Tag}}
