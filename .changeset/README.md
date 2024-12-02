# Changeset

Changeset is a tool used for managing changes to a project. It is a simple way to keep track of changes to a project and
to ensure that changes are properly documented. To use this, the first thing you need to do is install the command:

```bash
go install github.com/desmos-labs/changeset/cmd/changeset@latest
```

Then, each time you perform a change you can run the following command to create a new changeset:

```bash
changeset add
```

This will create a new changeset file in the `.changeset` directory that you can then commit to your repository.

Finally, when you are ready to release a new version of your project, you can run the following command:

```bash
changeset collect
```