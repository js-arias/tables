Tables
======

Tables is a simple tool form management of text based tables in the command
line.

It is based on the operator stream parading described in E. Schaffer & M.
Wolf, "[The UNIX shell as a fourth generation language](http://www.rdb.com/lib/4gl.pdf)"
UNIX Review, 8(3) (March 1991), p. 24. Under this paradigm an operator
(command) performs a unique function on the data.

Table format
------------

A table is an ordinary utf-8 file, in which each record (row) is delimited by
a newline character, and fields (columns) by tab (by default) character. The
first line of the file contains the header with the name of all columns, the
remain of the file contains the data.

An example dataset is:

	Item	Amount	Cost	Value	Description
	1	3	50	150	rubber gloves
	2	100	5	500	test tubes
	3	5	80	400	clamps
	4	23	19	437	plates
	5	99	24	2376	cleaning cloth
	6	89	147	13083	bunsen burners
	7	5	175	875	scales

Other similar (and more complete) tools
---------------------------------------

[/RBD](http://www.rdb.com/) is a commercial package, from the Revolution
Software which was used as the basis of all other packages.

[NoSQL](http://www.strozzi.it/cgi-bin/CSA/tw7/I/en_US/nosql/Home%20Page) by
Carlo Strozzi.

[Starbase](http://hopper.si.edu/wiki/mmti/Starbase) by John Roll optimized for
astronomic instrument data.


Authorship and license
----------------------

Copyright (c) 2016, J. Salvador Arias <jsalarias@csnat.unt.edu.ar>
All rights reserved.
Distributed under BSD2 license that can be found in the LICENSE file.

