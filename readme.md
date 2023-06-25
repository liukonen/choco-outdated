# choco-outdated - Faster Alternative to Chocolatey's "choco outdated" Command

This Go application is designed to provide a faster alternative to the Chocolatey application's `choco outdated` command. It checks for outdated packages and generates a report indicating the current and available versions of each package.

## Dependencies

This application uses the following dependencies:

- [github.com/spf13/viper](https://github.com/spf13/viper) - for reading the configuration file.

## Installation

To install the application, follow these steps:

1. Clone the repository:
```bash
   git clone <repository_url>
```
 2. Change to the application's directory:
``` bash
cd <application_directory>
```
3. Build the application:
``` bash
go build
```
4. Run the application:
``` bash
./<application_name>
```

## Configuration
Before running the application, you need to configure it by providing the necessary settings. Follow the steps below to configure the application:

1. Create a configuration file named config.yaml in the application's directory.

2. Specify the URLs to check for the latest versions of packages. Add the URLs under the urls field in the configuration file. For example:

``` yaml
urls:
  - https://example.com/packages
  - https://another-example.com/packages
```
3. In the same file, list out the location of where your choco file nespec files are installed. by default they exist in C:\ProgramData\chocolatey\lib\*\*.nuspec

## Output
When you run the application, it will display the following information for each package:

``` mathematica
Package Name | Current Version | Available Version
```

- If the available version is the same as the current version, it will be displayed in green.
- If the available version is different from the current version, it will be displayed in yellow.
- If the available version is unknown or unavailable, it will be displayed in red.
At the end of the output, if there are any available updates, the application will provide an install script to upgrade the packages using Chocolatey. You can run the provided install script to update the packages.

Please note that you need to have Chocolatey installed on your system for the install script to work.

## Contributing
If you find any issues or have suggestions for improvement, please feel free to contribute to this project. You can submit bug reports, feature requests, or pull requests through the project's repository.

## License
This application is licensed under the [MIT License.](LICENSE)
