import yaml
import pathlib

def main():
  # get all yaml files in PR (not just new/modified)
  yml_files = [p for p in pathlib.Path(".").rglob('*') if p.suffix in [".yml", ".yaml"]]
  # remove any under a secret directory (wha?)
  yml_files = [y for y in yml_files if "secret/" not in str(y)]


  # try to safe parse yaml files
  malformed_yaml = []
  for y in yml_files:
    with open(y) as stream:
      try:
        yaml.safe_load(stream)
      except yaml.YAMLError as exc:
        malformed_yaml.append((y, exc))
  
  print(malformed_yaml[0][0], malformed_yaml[0][1])


if __name__ == "__main__":
    main()
