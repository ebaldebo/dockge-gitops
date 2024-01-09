package git

const (
	repoUpToDateMsg         = "Repo is up to date"
	repoNotExistsCloningMsg = "Repo does not exist, cloning..."
	repoClonedMsg           = "Repo cloned"
	repoNotUpToDateMsg      = "Repo is not up to date, pulling..."
	repoPulledMsg           = "Repo pulled"
	copyingFilesMsg         = "Copying files..."
	filesCopiedMsg          = "Files copied"

	urlParseErr                     = "error parsing url: %w"
	gitFetchErr                     = "error fetching repo: %w"
	getLocalHashErr                 = "error getting local commit hash: %w"
	getRemoteHashErr                = "error getting remote commit hash: %w"
	checkingIfRepoHasUpdateErr      = "error checking if repo has update: %w"
	checkingIfRepoExistsErr         = "error checking if repo exists: %w"
	cloningRepoErr                  = "error cloning repo: %w"
	pullingRepoErr                  = "error pulling repo: %w"
	readingDirErr                   = "error reading dir: %w"
	cloneDirNotEmptyErr             = "error cloning into dir, dir is not empty: %w"
	gettingFilesFromDestinationErr  = "error getting files from destination: %w"
	removingFilesFromDestinationErr = "error removing files from destination: %w"
	copyingEnvFileErr               = "error copying .env file to desination dir: %w"
	copyingSubfoldersErr            = "error copying subfolders to desination dir: %w"
	gettingSubDirsError             = "error getting subdirs: %w"

	envFilePath = "/env/.env"
)