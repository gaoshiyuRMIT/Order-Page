### General Setup
#### Edit configuration files
- Specify database connection details in `config.json` (refer to `config.example.json`)

### API
Make sure you are in `app`: `cd app`

#### Setup
- `go mod download`

#### Usage
- `go run server.go`


### Data Import
Make sure you are in `import_script`: `cd import_script`

#### Set up the virtual environment
- create venv: `python3 -m venv venv`
- activate venv: `. venv/bin/activate`
- install libraries: `pip install -r requirements.txt`

#### Usage
- activate venv: - `. venv/bin/activate`
- Make sure the `test_data` folder, `import_data.py` and `config.json` are in the same directory.
- `python import_data.py`