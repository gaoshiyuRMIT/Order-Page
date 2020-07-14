#### Set up the virtual environment
- `python3 -m venv venv`
- `. venv/bin/activate`
- `pip install -r requirements.txt`

#### Edit configuration files
- Specify database connection details in `config.json` (refer to `config.example.json`)

#### Usage
- `. venv/bin/activate`
- Make sure the `test_data` folder, `import_data.py` and `config.json` are in the same directory.
- `python import_data.py`