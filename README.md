# Bloodhound CE JSON Uploader

Simple binary to upload files over the API to BloodHound CE.

## Usage

```bash
./BHCEupload -tokenid <tokenid> -tokenkey <tokenkey> -dir <dir> -url <url>
```

By default, the url is set to localhost:8080, and the dir is set to the current directory.
For large environments, with collections above 15-20Gigs, it is recommended to split the files for upload into smaller chunks before uploading otherwise the server will time out.

Additionally, make sure to have enough RAM or SWAP space to handle the large files. Since the uploads need to be signed they'll need to be loaded into memory before being sent to the server.

```markdown
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣴⣶⣾⠿⠶⢶⣶⣤⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⡿⠟⠋⠉⠀⠀⠀⠀⠀⠀⠈⠉⠉⠻⣧⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⠟⠉⠀⠀⣠⣤⣶⣷⢀⡀⠀⠀⠀⠀⠀⠀⠈⠻⣦⡀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢀⣼⡟⠁⠀⠀⠀⠘⠛⠁⠀⣀⣼⣿⣦⠀⠀⠀⠀⠀⠙⣦⣮⣹⡆⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣠⣾⡿⠀⠀⠀⠀⠀⠀⠀⠰⢿⡿⠿⠟⠛⠉⠀⠀⠀⠀⠀⠈⠙⠛⠳⠶⣤⣄⡀⠀⠀
⠀⠀⠀⠀⢀⣾⡿⠛⢀⡀⠀⠀⢀⣠⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡤⠚⣉⣉⣉⣙⣧⠀
⠀⠀⠀⠀⣼⠏⠀⠀⠀⣷⠀⣴⡟⠁⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⠀⣠⠀⠀⠀⣿⠀⢿⣿⣿⣿⣿⣿⠀
⠀⠀⠀⢀⡿⠀⠀⠀⠀⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⢀⣠⣴⣯⣤⡾⠃⠀⠀⠀⠈⠻⠶⣾⣿⣿⣿⣿⠄
⠀⠀⠀⣸⡇⠀⠀⠀⢰⣿⣿⠀⡄⠀⠀⠀⠀⠀⠀⠀⠀⠘⠛⠉⠀⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⡿⠀
⠀⠀⣰⡿⠀⠀⠀⠀⠈⢹⣿⡼⠁⠀⠀⠀⠀⠀⢠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣾⣿⣿⠟⠁⠀
⠀⣰⡿⢡⠂⠀⠀⠀⠀⠀⣿⣇⣰⣂⣠⠀⠀⠀⣟⣀⣀⠀⠀⠀⠀⠀⣠⣤⠶⠿⠿⠿⠛⠿⣧⠀⠀⠀
⣰⡿⢠⣏⡔⠀⠀⠀⠀⠀⢿⣿⣿⣿⡇⢠⡀⠀⢿⣿⣿⣿⣍⣙⠛⢻⣏⣤⠾⣦⡀⠀⠀⠀⢹⡇⠀⠀
⣿⣷⣿⣿⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⣧⠀⠈⠻⢿⣿⣿⣿⣿⣿⣿⡿⣶⣾⣿⣶⣦⣴⡟⠁⠀⠀
⠹⢿⣿⣿⡇⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⡏⠀⠀⠀⠀⠈⠉⠉⠉⠙⠿⠟⠁⠈⠙⠛⢻⡿⠀⠀⠀⠀
⠀⠀⠙⠻⣧⡀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⠋⠀⠀⡴⣾⣿⣷⣶⣤⣤⣤⣤⣤⣶⠤⠶⠛⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠙⢿⣶⣤⣤⣤⣴⣿⣿⣿⣿⠃⠀⠀⠀⢰⡟⢿⣿⣿⡿⠻⣿⣿⡿⠁⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠉⠙⠛⠛⠋⠁⠙⢿⣿⠀⢠⠄⠀⡼⠀⢸⣿⢹⡇⠀⢹⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣶⡟⣸⠁⠃⠀⣿⣿⡀⠧⠀⣸⣿⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣟⣴⣿⠀⠀⠀⣿⠹⡇⠀⢠⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⠛⣻⡿⢿⣆⢠⣿⠀⠁⢀⣾⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠾⠋⠁⢈⣿⣿⣿⡇⠀⣾⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⡿⠋⠉⢿⣾⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠋⠀⠀⠀⠈⠿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
BloodhoundCE json uploader⢀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
Please provide all required flags: -tokenid, -tokenkey
OPTIONAL: -url, -dir, -h  for help
```