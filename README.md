# grandctl

Let's start a golang command line project.

1. ~~lets use https://github.com/spf13/cobra for our cli framework~~
2. ~~lets use dep for dependency management~~
3. ~~lets read a simple ansible hosts file~~
4. ~~lets enable local command exec~~
5. ~~lets start off, running only local operations from the `master` node or `boot` node~~
- [x] lets level all remote operations to `terraform`
- [x] lets keep it simple, and support only INCEPTION `uninstall`, `install`, `hosts` operations
- [ ] lets also support construction of the ansible hosts file, just to see how hard it is to build that logic

# notes

1. add a new package dependency to gopkg.toml, after you have imported it
   `dep ensure -add github.com/spf13/cobra/cobra`

# examples

1. uninstall INCEPTION
   ```
   grandctl uninstall --gate stable
   ```

2. install INCEPTION
   ```
   grandctl install --gate stable
   ```

3. dump `config.yaml` and update from `~/.grandctl/config.yaml`
   ```
   grandctl init --gate stable 
   ```

4. build and deploy
   ```
   DEPLOY_USER=user DEPLOY_TARGET=9.x.x.x make deploy
   ```


