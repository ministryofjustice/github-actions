import yaml
import pathlib

def main():

  yml_files = [p for p in pathlib.Path(".").rglob('*') if p.suffix in [".yml", ".yaml"]]
  yml_files = [y for y in yml_files if "secret/" not in str(y)]

  # print(yml_files)

  malformed_yaml = []
  for y in yml_files:
    with open(y) as stream:
      try:
        yaml.safe_load(stream)
      except yaml.YAMLError as exc:
        malformed_yaml.append(f"\n{str(y)}:\n{str(exc)}")

  if malformed_yaml != []:
    error_message = "Malformed YAML detected:\n" + "\n".join(malformed_yaml)
    raise Exception(error_message)
  else:
    print("All YAML files OK!")

if __name__ == "__main__":
    main()
