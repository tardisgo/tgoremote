#build the library to link against normal go
haxe -main tardis.Go -cp tardis -dce full -D static_link -cpp tardis/cpplib
