# Change Log

## Unreleased



## 1.3.0 - 2016-02-24

### Added
 - Version flag (-v)
 - broadcastAllDsViewers for viewer count messaging


## 1.2.0 - 2016-02-23

### Changed
 - Fixed logging messages for socket handlers

### Added
 - Viewer count to map.html template
 - send json through websocket with viewer count and instructions to update layer
 - Touch screen support for drawing features (https://github.com/michaelguild13/Leaflet.draw) 

### Removed
 - Redundant logging from socketHandlers.go


## 1.1.0 - 2016-02-22

### Added
 - Added Socket hub

Initial Release

### Added
 - Web interface for drawing.
 - Go based RESTful GeoJson server.
 - Bolt database
 - Cache layer for database
