import os
import re

from comment_parser import comment_parser


class Extract:

  def __init__(self,
               content,
               token_name,
               token_package,
               file_name,
               line_number):

    self.content = content
    self.token_name = token_name
    self.token_package = token_package
    self.file_name = file_name
    self.line_number = line_number

  def __str__(self):
    r_content = repr(self.content)
    r_token = repr(self.token_name)
    r_file_name = repr(self.file_name)
    r_line_number = repr(self.line_number)

    return f'Extract({r_content}, {r_token}, {r_file_name}, {r_line_number})'

  def __repr__(self):
    return str(self)


stdlib_method_signature = re.compile(
  'func (.+)\([\w]+ context\.Context, [\w]+ \.\.\.core\.Value\) \(core\.Value, error\)'
)


def get_doc_extracts(go_file, directory, token_impl_map):
  go_lines, go_code = {}, None

  # Read the contents of a file as a string and also index lines by line
  # numbers. This is required later on to process string by lines.
  with open(go_file, 'r') as file:
    for no, line in enumerate(file.readlines()):
      go_lines[no] = line

    file.seek(0)
    go_code = file.read()

  # Extract all comments in a go source file. Each line in the comment is
  # returned as Comment instance, the comments are thus, not grouped.
  comments = comment_parser.extract_comments_from_str(
    go_code,
    mime='text/x-go',
  )

  comment_groups = []
  current_line = None

  # Group the comment lines as block of comments. Lines that have consecutive
  # line numbers are assumed to belong to the same comment block.
  for comment in comments:
    new_line = comment.line_number()

    if current_line is None or new_line != current_line + 1:
      comment_groups += [[]]

    comment_groups[-1].append(comment)
    current_line = new_line

  doc_groups = []

  # Get the entity that a comment block is talking about and check if
  # it belongs to an entity that we are interested in and aggregate
  # them.
  for group in comment_groups:
    group_end = group[-1].line_number()
    comment_related_to = go_lines.get(group_end)

    if not comment_related_to:
      continue

    match = stdlib_method_signature.search(
      comment_related_to,
    )

    if match is None:
      continue

    doc_groups.append((group, match))

  # Collate all the gathered information about the extract and build
  # Extracts.
  for (group, match) in doc_groups:
    content = list(map(
      lambda comment: comment.text().strip(),
      group,
    ))

    file_name = os.path.relpath(go_file, directory)
    token_package = os.path.dirname(file_name)

    impl_name = match.group(1)
    token_name = token_impl_map.get(impl_name, None)

    if token_name is None:
      token_name = impl_name
      print(f'Token against {impl_name} was not found')

    # TODO: token_name should be the name the method registers itself with
    # instead of the name of the method that implements it.
    yield Extract(
      content=content,
      token_name=token_name,
      token_package=token_package,
      line_number=group[-1].line_number(),
      file_name=file_name,
    )


def get_token_impl_map(go_file):
  # Extract all implementations and token names from indices.

  regex = re.compile(
    '"(?P<token>[A-Z0-9_]+)":[ \t]*(?P<impl>[\w\d]+)'
  )

  with open(go_file, 'r') as file:
    for match in regex.finditer(file.read()):
      yield { match.group('impl'): match.group('token') }
