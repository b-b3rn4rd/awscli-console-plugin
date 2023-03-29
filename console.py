import logging
import subprocess
import os
import sys
# Kludge because this will throw up when using AWS CLI v2 for unknown reasons (the file is right there,
# and the path also seems correct)
try:
    import site
    found_site = True
except:
    found_site = False

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

        # miserable hack because there are issue with AWS CLI V2 which causes this plugin to not find site.py
        if found_site:
            bin = os.path.join(site.USER_BASE, 'bin', 'awscli-console-plugin')
            if not os.path.isfile(bin):
                bin = os.path.join(sys.prefix, 'bin', 'awscli-console-plugin')
        else:
            bin = os.path.join(sys.prefix, 'bin', 'awscli-console-plugin')
            if not os.path.isfile(bin):
                # in addition when using pyenv, the Go executable will be installed somewhere among the shims
                # instead of a "normal" location
                if os.environ['PYENV_ROOT']:
                    bin = os.path.join(os.environ['PYENV_ROOT'], 'shims', 'awscli-console-plugin')

        cmd.append(bin)

        if parsed_globals.profile:
            cmd.append('--profile={}'.format(parsed_globals.profile))

        if options.output_only:
            cmd.append('--output')

        if options.timeout:
            cmd.append('--timeout={}'.format(options.timeout))

        res = subprocess.run(cmd, stdout=subprocess.PIPE, universal_newlines=True)

        print(res.stdout)
