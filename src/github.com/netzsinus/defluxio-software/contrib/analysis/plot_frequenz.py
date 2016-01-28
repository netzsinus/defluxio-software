# vim:fileencoding=utf-8
import matplotlib.pyplot as plt
import numpy as np
import datetime as dt

def load_data(filename):
  datafile = open(filename)
  datafile.readline() # skip the header
  data = np.loadtxt(datafile)
  time = [dt.datetime.fromtimestamp(ts) for ts in data[:,0]]
  return time, data[:,1]

rhrk_time, rhrk_freq = load_data("rhrk-frequenz.txt")
foo_time, foo_freq = load_data("foofrequenz.txt")
plt.title("Netzfrequenz: RHRK vs. Foo")
plt.xlabel("Zeit")
plt.ylabel("Frequenz [Hz]")
plt.plot(rhrk_time, rhrk_freq, 'r', label="RHRK")
plt.plot(foo_time, foo_freq, 'b', label="Foo")
plt.legend()

print "Mittelwert RHRK:", np.mean(rhrk_freq)
print "Mittelwert Foo:", np.mean(foo_freq)

#plt.figure()
#lastweek_time, lastweek_freq = load_data("lastweek-frequenz.txt")
#plt.title("Netzfrequenz: Letzte Woche")
#plt.xlabel("Zeit")
#plt.ylabel("Frequenz [Hz]")
#plt.plot(lastweek_time, lastweek_freq, 'b', label="Letzte Woche (ITWM)")
#plt.legend()



plt.show()
