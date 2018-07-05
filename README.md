# Qbackup

QB allows to backup files to the Amazon S3 cloud storage. QB minimize a storage space as it backups only a directory which has changes. Also, it allows specifying how many copies you want to store.

## How to backup to Amazon S3

Steps to follow:

  1. Create a config file

```

---
name: myBackup
description: Some test case
logfile: /var/log/backup/backup.log
home: /var/qsbackup  # Specify where to store a database (Optional value, default  /home/username/.config/qsbackup)
storage:
    type: aws
    aws_region: eu-west-1
    aws_bucket: myBackup
    aws_key: AWS_KEY
    aws_secret: AWS_SECRET
# =================================================================
# Directories
# =================================================================
dirs:
  - name: photo
    path: /home/username/Photos

  - name: docs
    path: /home/username/Documents
```

  2. Run the program

**Extra commands:**

Last backups: `/usr/local/sbin/qsbackup -c /usr/local/etc/qsbackup.conf -l`

All backups: `/usr/local/sbin/qsbackup -c /usr/local/etc/qsbackup.conf -s`