#include <stdio.h>
#include <string.h>
#include hashMap.h

/* Given a feat file in standard format, constructs the 
 * lexicon hashmap and returns it. */

hash *hashify(char *filename) {
	FILE *file = fopen(filename, "rw");
	if (!file) {
		return NULL;
	}
	
}

