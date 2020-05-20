import re


def de_parenthesis(text):
  if text.startswith('('):
    text = text[1:]

  if text.endswith(')'):
    text = text[:-1]

  return text


def preprocess_line(target):
  target = target.replace('@params', '@param')
  target = target.replace('@returns', '@return')

  return target


def process_decorated(regex, target):
  match = regex.search(target)

  arg_name = None
  arg_type = match.group('type')
  arg_desc = match.group('desc')

  try:
    arg_name = match.group('name')
  except:
    pass

  arg_type = de_parenthesis(
    arg_type,
  )

  arg_types = list(map(
    str.strip,
    filter(None, arg_type.split('|')),
  ))

  arg_desc = arg_desc.capitalize()

  return arg_name, arg_types, arg_desc


def process_returns(returns):
  if returns is None:
    return None

  regex = re.compile(
    '@return (?P<type>\(.+\)) - (?P<desc>.+)',
  )

  arg_name, arg_types, arg_desc = process_decorated(
    regex,
    returns,
  )

  return {
    'name': arg_name,
    'type': arg_types,
    'desc': arg_desc,
  }


def process_params(params):
  regex = re.compile(
    '@param (?P<name>.+) (?P<type>\(.+\)) - (?P<desc>.+)',
  )

  for param in params:
    arg_name, arg_types, arg_desc = process_decorated(
      regex,
      param,
    )

    yield {
      'name': arg_name,
      'type': arg_types,
      'desc': arg_desc,
    }


def transform_extract(extract):
  content, params, returns = [], [], None

  for line in extract.content:
    line = preprocess_line(line.strip())

    if line.startswith('@param'):
      params.append(line)
      continue

    if line.startswith('@return'):
      returns = line
      continue

    if returns is not None:
      returns = f'{returns} {line}'
      continue

    if len(params) > 0:
      params[-1] = f'{params[-1]} {line}'
      continue

    if len(params) == 0:
      content.append(line)
      continue

  content = ' '.join(content)
  returns = process_returns(returns)
  params = list(process_params(params))

  return {
    'desc': content,
    'params': params,
    'returns': returns,
  }
