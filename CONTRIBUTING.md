# Contributing

To deploy a new release, tag the new release with a version beginning with `v`.

```sh
git tag v1
```

After tagging the release, push to the remote.

```sh
git push --follow-tags
```

After a new version is pushed, it will be automatically deployed.
