# Fare Estimator

The script is consisted of the following parts

- `Estimator`: Estimator is responsible to coordinate the whole process, it is creating the channels, files, writers and calls the go routines.
- `Scanner`: Scanner is responsible for the hard job of fetching the ride points data from the initial csv file
- `Worker`: It's job is to calculate the distance, time interval and speed and then export the the total ride fare.
- `Writer`: Writer is receives the final data from a channel and then exports to a file

In order fo the script to work, please replace the existing paths.csv file to the working directory