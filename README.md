# condenser

Podcast search & normalisation service.

Leverages the iTunes API to search for podcasts, extracts track data from the podcast feed to be supplied to a data store.

## running

Uses [glide](https://github.com/Masterminds/glide) for package management.

```bash
git clone https://github.com/jimah/condenser $GOPATH/src/condenser
cd $GOPATH/src/condenser
glide install
make dev
```