import logging
import subprocess
import os
import sys
import site
from awscli.customizations.commands import BasicCommand


def awscli_initialize(cli):
    cli.register('building-command-table.main', inject_commands)


def inject_commands(command_table, session, **kwargs):
    command_table['console'] = Console(session)


class Console(BasicCommand):
    NAME = 'console'
    DESCRIPTION = 'Authenticate to AWS console'
    SYNOPSIS = 'aws console [--profile=Name] [--timeout=Timeout] [--output-only=true|false]'

    ARG_TABLE = [
        {
            'name': 'timeout',
            'default': '',
            'help_text': 'Console session timeout in seconds, only for IAM user credentials'
        },
        {
            'name': 'output-only',
            'cli_type_name': 'boolean',
            'default': False,
            'help_text': 'Print the console url instead of opening it in the browser'
        },
    ]

    UPDATE = False

    def _run_main(self, args, parsed_globals):
        """Run the command and report success."""
        logging.basicConfig(level=logging.INFO)
        for handler in logging.root.handlers:
            handler.addFilter(logging.Filter(__name__))
        self._call(args, parsed_globals)

        return 0

    def _call(self, options, parsed_globals):
        """Run the command."""
        cmd = []

        bin = os.path.join(site.USER_BASE, 'bin', 'console')
        if not os.path.isfile(bin):
            bin = os.path.join(sys.prefix, 'bin', 'console')

        cmd.append(bin)

        if parsed_globals.profile:
            cmd.append('--profile={}'.format(parsed_globals.profile))

        if options.output_only:
            cmd.append('--output')

        if options.timeout:
            cmd.append('--timeout={}'.format(options.timeout))

        res = subprocess.run(cmd, stdout=subprocess.PIPE, universal_newlines=True)

        print(res.stdout)
