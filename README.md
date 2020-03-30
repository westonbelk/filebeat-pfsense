# Filebeat for pfSense

This project maintains a custom filebeat module and processor to parse, enrinch, and ship logs from pfSense directly to Elasticsearch.

## Features
- No need for a Logstash installation.
- Uses the standard syslog shipping implementation in pfSense. (No need for a separate package)
- Enrichment of firewall logs to include the name and description of the rules.

## Installation

Unfortunately, the Elastic Beats project does not provide builds for FreeBSD so Filebeat must be built from source. Additionally, this project implements a custom processor that we need to inject into the build. There are plugins in the Elastic Beats project to add custom processors to an existing installation, but this is not currently supported on the FreeBSD platform. 


## Issues (Todo)

- Fails to grok IGMP packets
- Need to clean up field organization
- Add support for more than just filterlog type
- Add meta information for module and processor
- Add installation instructions
- Create build script