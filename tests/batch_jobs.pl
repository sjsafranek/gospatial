#!/usr/bin/perl 

use strict;
use warnings;
use 5.012;
use File::Glob;
use Text::CSV;

sub trim { 
	# Lacey Powers
	my $string = shift;
    $string =~ s/^\s+//;
    $string =~ s/\s+$//;
    $string =~ s/^\t+//g;
    $string =~ s/\t$//g;
    $string =~ s/^\n//;
    $string =~ s/\n$//;
    $string =~ s/^\r//;
    $string =~ s/\r$//;
    $string =~ s/"//g;
    $string =~ s/\$//g;
    $string =~ s/^\.00$/0/g;
    # This has to be the last substitution for proper spacing.
    $string =~ s/\s{2,}/ /g;
    # http://perlmaven.com/trim
	$string =~ s/^\s+|\s+$//g; 
	return $string 
};

my $file = "service-perfdata.out";
print("$file\n");

open(my $fh, '<:encoding(UTF-8)', $file)
	or die "Could not open file '$file' $!";

my %services = ();

while (my $row = <$fh>) {
	chomp($row);

	# Service data
	# $LASTSERVICECHECK$\t$HOSTNAME$\t$SERVICEDESC$\t$SERVICESTATE$\t$SERVICEATTEMPT$\t$SERVICESTATETYPE$\t$SERVICEEXECUTIONTIME$\t$SERVICELATENCY$\t$SERVICEOUTPUT$\t$SERVICEPERFDATA$\n
	# https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/4/en/perfdata.html

	my @fields = split("\t", $row);
	my $timestamp = $fields[0];
	my $server = $fields[1];
	my $service = $fields[2];
	# my $status = $fields[3];
	# my $attempt = $fields[4];
	# my $type = $fields[5];
	my $executeTime = $fields[6];
	my $latency = $fields[7];

	$service =~ s#/##g;

	my $filename = "./static/nagios_service__${service}.csv";
	
	if (!-f $filename) {
		open(my $save_fh, '>>', $filename) 
			or die "Could not open file '$filename' $!";
		# print $save_fh "LASTSERVICECHECK\tHOSTNAME\tSERVICEDESC\tSERVICESTATE\tSERVICEATTEMPT\tSERVICESTATETYPE\tSERVICEEXECUTIONTIME\tSERVICELATENCY\n";
		print $save_fh "LASTSERVICECHECK\tHOSTNAME\tSERVICEDESC\tSERVICEEXECUTIONTIME\tSERVICELATENCY\n";
		# print $save_fh "$timestamp\t$server\t$service\t$status\t$attempt\t$type\t$executeTime\t$latency\n";
		print $save_fh "$timestamp\t$server\t$service\t$executeTime\t$latency\n";
		close $save_fh;
	}
	else {
		open(my $save_fh, '>>', $filename) 
			or die "Could not open file '$filename' $!";
		# print $save_fh "$timestamp\t$server\t$service\t$status\t$attempt\t$type\t$executeTime\t$latency\n";
		print $save_fh "$timestamp\t$server\t$service\t$executeTime\t$latency\n";
		close $save_fh;
	}
}


# Close files
# foreach my $service ( keys %services ) {
# 	close $services{$service}
# }

close $fh;


# Host data
# $LASTHOSTCHECK$\t$HOSTNAME$\t$HOSTSTATE$\t$HOSTATTEMPT$\t$HOSTSTATETYPE$\t$HOSTEXECUTIONTIME$\t$HOSTOUTPUT$\t$HOSTPERFDATA$\n 
# 
# Service data
# $LASTSERVICECHECK$\t$HOSTNAME$\t$SERVICEDESC$\t$SERVICESTATE$\t$SERVICEATTEMPT$\t$SERVICESTATETYPE$\t$SERVICEEXECUTIONTIME$\t$SERVICELATENCY$\t$SERVICEOUTPUT$\t$SERVICEPERFDATA$\n
# https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/4/en/perfdata.html
