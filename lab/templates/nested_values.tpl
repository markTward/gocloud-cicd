service:
  gocloudAPI:
    image:
      repository: {{.Dynamic.Repo}}
      tag: {{.Dynamic.Tag}}
  gocloudGrpc:
    image:
      repository: {{.Dynamic.Repo}}
      tag: {{.Dynamic.Tag}}
  abc:
    debug: {{.A.B.C}}
